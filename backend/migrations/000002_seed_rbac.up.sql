INSERT INTO roles (name, description)
VALUES
    ('admin', 'Full system access'),
    ('manager', 'Can manage users and view audit logs'),
    ('analyst', 'Can view users and audit logs'),
    ('viewer', 'Read-only access')
    ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (key, description)
VALUES
    ('users:read', 'View users'),
    ('users:write', 'Create and update users'),
    ('roles:read', 'View roles'),
    ('roles:write', 'Create and update roles'),
    ('audit_logs:read', 'View audit logs'),
    ('settings:update', 'Update system settings')
    ON CONFLICT (key) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         CROSS JOIN permissions p
WHERE r.name = 'admin'
    ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p ON p.key IN (
                                         'users:read',
                                         'users:write',
                                         'roles:read',
                                         'audit_logs:read'
    )
WHERE r.name = 'manager'
    ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p ON p.key IN (
                                         'users:read',
                                         'audit_logs:read'
    )
WHERE r.name = 'analyst'
    ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p ON p.key IN (
    'users:read'
    )
WHERE r.name = 'viewer'
    ON CONFLICT DO NOTHING;