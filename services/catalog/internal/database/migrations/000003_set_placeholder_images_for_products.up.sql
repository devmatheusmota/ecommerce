-- Set a placeholder image for all products that have no images
UPDATE products
SET images = '["/placeholder-product.svg"]'::jsonb
WHERE images = '[]'::jsonb OR images IS NULL;
