CREATE OR REPLACE VIEW ITEMS_BY_OBJECT_VIEW AS
SELECT
    t_item.*,
    t_link.OBJECT_ID AS OBJECT_ID
FROM ITEMS t_item
INNER JOIN OBJECT_ITEM_LINK t_link ON t_item.ID = t_link.ITEM_ID;
