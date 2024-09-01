# Prepare database.
mongo \
    -u ${MONGO_INITDB_ROOT_USERNAME} \
    -p ${MONGO_INITDB_ROOT_PASSWORD} \
    --authenticationDatabase admin ${MONGO_NAME} \
<<-EOJS
use ${MONGO_NAME};
db.createCollection("report");
EOJS

mongo \
    -u ${MONGO_INITDB_ROOT_USERNAME} \
    -p ${MONGO_INITDB_ROOT_PASSWORD} \
    --authenticationDatabase "$rootAuthDatabase" ${MONGO_NAME} \
<<-EOJS
db.createUser({
    user: "${MONGO_USER}",
    pwd: "${MONGO_PASS}",
    roles: [{
        role: "readWrite",
        db: "${MONGO_NAME}"
    }],
    mechanisms: ["${MONGO_AUTH}"],
});
EOJS
 
