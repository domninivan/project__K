package ru.mirea.production.security;

import java.util.*;
import java.time.Instant;
import javax.sql.DataSource;
import com.nimbusds.jose.*;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jwt.*;
import org.junit.jupiter.api.*;
import org.springframework.beans.factory.xml.XmlBeanDefinitionReader;
import org.springframework.context.support.PropertySourcesPlaceholderConfigurer;
import org.springframework.core.env.MapPropertySource;
import org.springframework.jdbc.datasource.DriverManagerDataSource;
import org.springframework.mock.web.MockServletContext;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.web.FilterChainProxy;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.web.context.support.GenericWebApplicationContext;
import org.springframework.web.bind.annotation.*;
import ru.mirea.production.user.entity.Role;

import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

class SecurityXmlTest {
    private static final byte[] KEY = new byte[32];
    private GenericWebApplicationContext context;
    private Accounts accounts;
    private TokenService tokens;
    private MockMvc mvc;

    @BeforeEach
    void setup() {
        context = new GenericWebApplicationContext(new MockServletContext());
        context.getEnvironment().getPropertySources().addFirst(new MapPropertySource("test", Map.of(
                "JWT_SECRET_BASE64", Base64.getEncoder().encodeToString(KEY))));
        var placeholders = new PropertySourcesPlaceholderConfigurer();
        placeholders.setEnvironment(context.getEnvironment());
        context.addBeanFactoryPostProcessor(placeholders);
        new XmlBeanDefinitionReader(context).loadBeanDefinitions("classpath:security-context.xml");
        context.removeBeanDefinition("accountDetailsService");
        accounts = new Accounts();
        context.getBeanFactory().registerSingleton("accountDetailsService", accounts);
        context.refresh();
        tokens = context.getBean(TokenService.class);
        var controller = new AuthController(context.getBean("authenticationManager", AuthenticationManager.class), tokens);
        mvc = MockMvcBuilders.standaloneSetup(controller, new ProbeController())
                .addFilters(context.getBean(FilterChainProxy.class)).build();
    }

    @AfterEach
    void close() { context.close(); }

    @Test
    void loginMeAndNoPasswordDisclosure() throws Exception {
        var response = mvc.perform(post("/api/auth/login").contentType("application/json")
                        .content("{\"username\":\"client\",\"password\":\"test-password\"}"))
                .andExpect(status().isOk()).andExpect(jsonPath("$.tokenType").value("Bearer"))
                .andExpect(header().string("Cache-Control", "no-store")).andReturn();
        var json = new com.fasterxml.jackson.databind.ObjectMapper().readTree(response.getResponse().getContentAsString());
        mvc.perform(get("/api/me").header("Authorization", "Bearer " + json.get("accessToken").asText()))
                .andExpect(status().isOk()).andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.role").value("CLIENT"))
                .andExpect(jsonPath("$.password").doesNotExist()).andExpect(jsonPath("$.passwordHash").doesNotExist());
    }

    @Test
    void wrongUnknownAndDisabledCredentialsReturnSameError() throws Exception {
        for (String username : List.of("client", "missing")) {
            mvc.perform(post("/api/auth/login").contentType("application/json")
                            .content("{\"username\":\"" + username + "\",\"password\":\"wrong\"}"))
                    .andExpect(status().isUnauthorized()).andExpect(jsonPath("$.code").value("UNAUTHORIZED"));
        }
        accounts.active = false;
        mvc.perform(post("/api/auth/login").contentType("application/json")
                        .content("{\"username\":\"client\",\"password\":\"test-password\"}"))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void missingMalformedExpiredAndWrongSignatureAreRejected() throws Exception {
        mvc.perform(get("/api/me")).andExpect(status().isUnauthorized());
        for (String token : List.of("garbage", signed(KEY, Instant.now().minusSeconds(300)),
                signed("different-signing-key-for-testing-only".getBytes(java.nio.charset.StandardCharsets.UTF_8),
                        Instant.now().plusSeconds(300)))) {
            mvc.perform(get("/api/me").header("Authorization", "Bearer " + token))
                    .andExpect(status().isUnauthorized());
        }
    }

    @Test
    void rolesAreEnforcedAndRefreshedFromDatabase() throws Exception {
        String token = tokens.issue(accounts.loadById(1L)).accessToken();
        mvc.perform(get("/api/planner/probe").header("Authorization", "Bearer " + token))
                .andExpect(status().isForbidden());
        accounts.role = Role.PLANNER;
        mvc.perform(get("/api/planner/probe").header("Authorization", "Bearer " + token))
                .andExpect(status().isOk());
        accounts.role = Role.REPORTING;
        mvc.perform(get("/api/products").header("Authorization", "Bearer " + token))
                .andExpect(status().isForbidden());
        accounts.active = false;
        mvc.perform(get("/api/me").header("Authorization", "Bearer " + token))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void corsAllowsConfiguredFrontendAndRejectsOtherOrigins() throws Exception {
        mvc.perform(options("/api/orders").header("Origin", "http://localhost:5173")
                        .header("Access-Control-Request-Method", "POST")
                        .header("Access-Control-Request-Headers", "Authorization,Content-Type"))
                .andExpect(status().isOk())
                .andExpect(header().string("Access-Control-Allow-Origin", "http://localhost:5173"));
        mvc.perform(options("/api/orders").header("Origin", "https://unknown.example")
                        .header("Access-Control-Request-Method", "POST"))
                .andExpect(status().isForbidden());
    }

    @Test
    void refusesWeakSigningKey() {
        assertThrows(IllegalArgumentException.class,
                () -> new TokenService(Base64.getEncoder().encodeToString(new byte[8]), 900, "case22-java"));
    }

    private String signed(byte[] key, Instant expires) throws Exception {
        var jwt = new SignedJWT(new JWSHeader(JWSAlgorithm.HS256), new JWTClaimsSet.Builder()
                .subject("1").issuer("case22-java").issueTime(Date.from(Instant.now().minusSeconds(600)))
                .expirationTime(Date.from(expires)).build());
        jwt.sign(new MACSigner(key));
        return jwt.serialize();
    }

    static class Accounts extends AccountDetailsService {
        boolean active = true;
        Role role = Role.CLIENT;
        final String hash = new BCryptPasswordEncoder().encode("test-password");
        Accounts() { super(new DriverManagerDataSource()); }
        @Override public AppPrincipal loadUserByUsername(String name) {
            if (!name.equals("client")) throw new UsernameNotFoundException("Missing");
            return new AppPrincipal(1L, "Client", "client", hash, role, active);
        }
        @Override public AppPrincipal loadById(Long id) {
            if (id != 1L) throw new UsernameNotFoundException("Missing");
            return loadUserByUsername("client");
        }
    }

    @RestController
    static class ProbeController {
        @GetMapping({"/api/planner/probe", "/api/products"})
        String probe() { return "ok"; }
    }
}
