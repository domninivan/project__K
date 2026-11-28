package ru.mirea.production.security;

import com.nimbusds.jose.jwk.source.ImmutableSecret;
import javax.crypto.spec.SecretKeySpec;
import java.time.Instant;
import java.util.Base64;
import org.springframework.security.oauth2.jose.jws.MacAlgorithm;
import org.springframework.security.oauth2.jwt.*;
import org.springframework.security.oauth2.core.*;

public class TokenService {
    private final JwtEncoder encoder;
    private final NimbusJwtDecoder decoder;
    private final long ttlSeconds;
    private final String issuer;

    public TokenService(String base64Secret, long ttlSeconds, String issuer) {
        byte[] bytes = Base64.getDecoder().decode(base64Secret);
        if (bytes.length < 32) throw new IllegalArgumentException("JWT secret requires at least 32 random bytes");
        if (ttlSeconds < 60 || ttlSeconds > 86400) throw new IllegalArgumentException("Invalid JWT TTL");
        var key = new SecretKeySpec(bytes, "HmacSHA256");
        this.encoder = new NimbusJwtEncoder(new ImmutableSecret<>(key));
        this.decoder = NimbusJwtDecoder.withSecretKey(key).macAlgorithm(MacAlgorithm.HS256).build();
        OAuth2TokenValidator<Jwt> requiredClaims = jwt -> jwt.getExpiresAt() != null && jwt.getIssuedAt() != null
                && jwt.getSubject() != null ? OAuth2TokenValidatorResult.success()
                : OAuth2TokenValidatorResult.failure(new OAuth2Error("invalid_token"));
        decoder.setJwtValidator(new DelegatingOAuth2TokenValidator<>(
                JwtValidators.createDefaultWithIssuer(issuer), requiredClaims));
        this.ttlSeconds = ttlSeconds;
        this.issuer = issuer;
    }

    public JwtDecoder decoder() { return decoder; }

    public TokenResponse issue(AppPrincipal user) {
        Instant now = Instant.now();
        var claims = JwtClaimsSet.builder().issuer(issuer).subject(user.getId().toString())
                .issuedAt(now).expiresAt(now.plusSeconds(ttlSeconds)).build();
        String token = encoder.encode(JwtEncoderParameters.from(
                JwsHeader.with(MacAlgorithm.HS256).build(), claims)).getTokenValue();
        return new TokenResponse(token, "Bearer", ttlSeconds);
    }

    public record TokenResponse(String accessToken, String tokenType, long expiresIn) {}
}
