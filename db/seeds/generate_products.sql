TRUNCATE TABLE products CASCADE;

INSERT INTO products (id, seller_id, name, price, stock, discount, type, description, created_at, updated_at)
SELECT 
    gen_random_uuid() AS id,
    '00000000-0000-0000-0000-000000000000'::uuid AS seller_id,
    CASE (i % 10)
        WHEN 0 THEN 'Bandai HG Gundam Aerial Skala 1/144 ' || i
        WHEN 1 THEN 'Action Figure Monkey D. Luffy Gear 5 Wano ' || i
        WHEN 2 THEN 'Lego Star Wars Ultimate Millennium Falcon ' || i
        WHEN 3 THEN 'Tamiya Mini 4WD Magnum Saber Premium ' || i
        WHEN 4 THEN 'Hasbro Transformers Optimus Prime Voyager Class ' || i
        WHEN 5 THEN 'Good Smile Nendoroid Hatsune Miku 2.0 ' || i
        WHEN 6 THEN 'Nintendo Switch Pokémon Violet Game Disc ' || i
        WHEN 7 THEN 'Hot Toys Iron Mark LXXXV Endgame Figure ' || i
        WHEN 8 THEN 'Boardgame Catan Standard English Edition ' || i
        ELSE 'Tamiya Spray Paint TS-14 Black Glossy ' || i
    END AS name,
    (50000 + (i * 12500) % 2500000) AS price,
    (5 + (i * 3) % 190) AS stock,
    CASE WHEN (i % 3 = 0) THEN (5 + (i * 5) % 45) ELSE 0 END AS discount,
    CASE (i % 10)
        WHEN 0 THEN 'model-kits'
        WHEN 1 THEN 'action-figure'
        WHEN 2 THEN 'collectibles'
        WHEN 3 THEN 'model-kits'
        WHEN 4 THEN 'action-figure'
        WHEN 5 THEN 'action-figure'
        WHEN 6 THEN 'games-puzzles'
        WHEN 7 THEN 'action-figure'
        WHEN 8 THEN 'games-puzzles'
        ELSE 'collectibles'
    END AS type,
    'Deskripsi premium untuk produk hobi ke-' || i || '. Mainan koleksi berkualitas tinggi untuk penggemar sejati. Kondisi mulus, original, dan siap kirim ke seluruh Indonesia.' AS description,
    NOW() - (i * INTERVAL '10 minute') AS created_at,
    NOW() - (i * INTERVAL '10 minute') AS updated_at
FROM generate_series(1, 1000) AS i;

