package ru.mirea.production.security;

import org.springframework.core.convert.converter.Converter;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.oauth2.core.*;
import org.springframework.security.oauth2.jwt.Jwt;

public class AccountJwtConverter implements Converter<Jwt, AbstractAuthenticationToken> {
    private final AccountDetailsService accounts;
    public AccountJwtConverter(AccountDetailsService accounts) { this.accounts = accounts; }

    @Override
    public AbstractAuthenticationToken convert(Jwt jwt) {
        try {
            AppPrincipal user = accounts.loadById(Long.valueOf(jwt.getSubject()));
            if (!user.isEnabled()) throw new UsernameNotFoundException("Disabled");
            return UsernamePasswordAuthenticationToken.authenticated(user, null, user.getAuthorities());
        } catch (NumberFormatException | UsernameNotFoundException ex) {
            throw new OAuth2AuthenticationException(new OAuth2Error("invalid_token"));
        }
    }
}
