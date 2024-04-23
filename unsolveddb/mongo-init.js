db.createUser(
    {
        user: "api",
        pwd: "password",
        roles: [
            {
                role: "readWrite",
                db: "unsolved-db"
            }
        ]
    }
)

db = new Mongo().getDB("unsolved-db")

db.createCollection("tickets")

db.tickets.insertOne(
    {
        tags: ["resolved", "management"],
        problem: "Test: Escalation to Management",
    }
)
