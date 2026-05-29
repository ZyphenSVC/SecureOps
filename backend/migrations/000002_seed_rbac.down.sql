DELETE FROM role_permissions;

DELETE FROM permissions
WHERE key IN (
    'users:read',
    'users:write',
    'roles:read',
    'roles:write',
    'audit_logs:read',
    'settings:update'
    );

DELETE FROM roles
WHERE name IN (
               'admin',
               'manager',
               'analyst',
               'viewer'
    );