-- Print agents report their version and operating system with every hello, so admins see what runs at the range.
ALTER TABLE print_agents ADD COLUMN version text NOT NULL DEFAULT '', ADD COLUMN os text NOT NULL DEFAULT '';
