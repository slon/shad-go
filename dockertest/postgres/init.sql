create table Role(
    RoleID SERIAL PRIMARY KEY,
    RoleName varchar(50)
);

insert into Role(RoleName)
values ('Admin'),('User');
