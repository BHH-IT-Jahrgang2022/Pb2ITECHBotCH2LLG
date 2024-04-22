db.createUser(
    {
        user: "api",
        pwd: "password",
        roles: [
            {
                role: "readWrite",
                db: "matcher-db"
            }
        ]
    }
)

db = new Mongo().getDB("matcher-db")
db.createCollection("windowfly")
db.createCollection("cleanbug")
db.createCollection("gardenbeetle")
db.createCollection("empty")

db.windowfly.insertOne(
    {
        keywords: ["Windowfly", "ab.*st(ü|u|ue)rz"],
        response: "Bitte wärmen und unter einer Decke schlafen lassen."
    }
)

db.cleanbug.insertOne(
    {
        keywords: ["Cleanbug", "an.*halt"],
        response: "Bitte tauchen Sie den Cleanbug unter Wasser, am Besten entweder in einer Badewanne oder in einem See."
    }
)

db.gardenbeetle.insertOne(
    {
        keywords: ["Gardenbeetle", "aus.*spuck"],
        response: "Bitte schockfrosten sie ihr Gerät."
    }
)

db.empty.insertOne(
    {
        keywords: [""],
        response: "Entschuldigen Sie bitte, ich habe Sie nicht verstanden. Könnten Sie das bitte weitere Informationen bereitstellen?\nAm Besten nennen Sie das betroffene Produkt und das konkrete Problem."
    }
)

db.empty.insertOne(
    {
        keywords: ["danke"],
        response: "Ich hoffe ich konnte Ihnen helfen. Falls Sie weitere Fragen haben, stehe ich Ihnen gerne zur Verfügung."
    }
)
