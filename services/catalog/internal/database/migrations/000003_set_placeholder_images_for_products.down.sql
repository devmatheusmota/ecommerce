-- Remove placeholder image from products that only have the default placeholder
UPDATE products
SET images = '[]'::jsonb
WHERE images = '["/placeholder-product.svg"]'::jsonb;
