package ru.mirea.production.product.dao;

import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;
import ru.mirea.production.product.entity.ProductMaterialRequirement;

import java.util.List;

@Repository
public class ProductMaterialRequirementDAOImpl implements ProductMaterialRequirementDAO{

    @Autowired
    private SessionFactory sessionFactory;

    @Override
    public List<ProductMaterialRequirement> getProductMaterialRequirementsByProductId(int productId){
        Session session = sessionFactory.getCurrentSession();
        return session
                .createQuery("from ProductMaterialRequirement where product.id = :product_id", ProductMaterialRequirement.class)
                .setParameter("product_id", productId)
                .getResultList();
    }
}
