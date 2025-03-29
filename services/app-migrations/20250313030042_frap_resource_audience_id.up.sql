
-- flattened resource_audience_policies_table
-- View to flatten the resource_audience for checking users' resource access easily
DROP VIEW IF EXISTS app.flattened_resource_audience_policies;
CREATE OR REPLACE VIEW app.flattened_resource_audience_policies AS
    -- Get users
    SELECT 
        ra.resource_id,
        ra.resource_audience_type,
        ra.resource_audience_id AS user_id,  -- Direct user access
        ra.resource_type,
        ra.privilege,
        ra.created_at,
        ra.updated_at,
        ra.deleted_at,
        ra.resource_audience_policy_id,
        ra.resource_audience_id AS resource_audience_id
    FROM 
        app.resource_audience_policies ra
    WHERE 
        ra.resource_audience_type = 'user'
    
    UNION ALL

    -- Flatten team members
    SELECT 
        ra.resource_id,
        ra.resource_audience_type,
        tm.user_id AS user_id,
        ra.resource_type,
        ra.privilege,
        ra.created_at,
        ra.updated_at,
        ra.deleted_at,
        ra.resource_audience_policy_id,
        ra.resource_audience_id AS resource_audience_id
    FROM 
        app.resource_audience_policies ra
    JOIN 
        app.team_memberships tm 
        ON ra.resource_audience_type = 'team' 
        AND ra.resource_audience_id = tm.team_id
        AND tm.deleted_at IS NULL
        AND ra.deleted_at IS NULL

    UNION ALL

    -- Flatten organization members
    SELECT 
        ra.resource_id,
        ra.resource_audience_type,
        jra.resource_audience_id AS user_id,
        ra.resource_type,
        ra.privilege,
        ra.created_at,
        ra.updated_at,
        ra.deleted_at,
        ra.resource_audience_policy_id,
        ra.resource_audience_id AS resource_audience_id
    FROM
        app.resource_audience_policies ra
    JOIN
        app.resource_audience_policies jra
        ON ra.resource_audience_type = 'organization'
        AND jra.resource_type = 'organization'
        AND jra.resource_id = ra.resource_audience_id
        AND jra.resource_audience_type = 'user'
        AND jra.deleted_at IS NULL