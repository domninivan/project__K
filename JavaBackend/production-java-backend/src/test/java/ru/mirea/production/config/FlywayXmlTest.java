package ru.mirea.production.config;

import org.flywaydb.core.Flyway;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.xml.XmlBeanDefinitionReader;
import org.springframework.context.support.GenericApplicationContext;
import org.springframework.jdbc.datasource.DriverManagerDataSource;

import javax.sql.DataSource;

import static org.junit.jupiter.api.Assertions.*;

class FlywayXmlTest {
    @Test
    void xmlCreatesFlywayOnExistingDataSourceBeforeHibernate() {
        try (var context = new GenericApplicationContext()) {
            new XmlBeanDefinitionReader(context).loadBeanDefinitions("classpath:applicationContext.xml");
            assertArrayEquals(new String[]{"javaFlyway"},
                    context.getBeanDefinition("sessionFactory").getDependsOn());
            assertEquals("migrate", context.getBeanDefinition("javaFlyway").getInitMethodName());

            // Verify XML wiring without running migrations or connecting to PostgreSQL.
            context.getBeanDefinition("javaFlyway").setInitMethodName(null);
            for (String bean : new String[]{"sessionFactory", "transactionManager", "goRestClient", "goRestClientBuilder", "dataSource"}) {
                context.removeBeanDefinition(bean);
            }
            DataSource dataSource = new DriverManagerDataSource();
            context.getBeanFactory().registerSingleton("dataSource", dataSource);
            context.refresh();

            var config = context.getBean("javaFlyway", Flyway.class).getConfiguration();
            assertSame(dataSource, config.getDataSource());
            assertEquals("public", config.getDefaultSchema());
            assertEquals("1", config.getBaselineVersion().toString());
            assertTrue(config.isCleanDisabled());
            assertEquals("classpath:db/migration", config.getLocations()[0].getDescriptor());
        }
    }
}
