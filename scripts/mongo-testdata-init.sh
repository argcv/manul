#!/usr/bin/env bash

mongo admin --eval 'db.createUser({ user: "manul", pwd: "pass4manul", roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "root","dbAdmin","dbOwner" ] });'

echo "Created User manul"
echo "db name: admin"
echo "user:    manul"
echo "pass:    pass4manul"
