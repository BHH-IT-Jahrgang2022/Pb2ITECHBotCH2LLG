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
        response: "Nach einem Absturz braucht die Windowfly Ruhe. Bitte wärmen Sie sie und lassen Sie sie unter einer Decke schlafen."
    }
)

db.windowfly.insertOne(
    {
        keywords: ["Windowfly", "l(ä|ae|a)dt.*nicht|nicht.*laden"],
        response: "Stellen Sie sicher, dass das Ladekabel korrekt angeschlossen ist und die Steckdose Strom führt. Überprüfen Sie auch, ob die Kontakte am Gerät sauber und frei von Schmutz sind."
    }
)

db.windowfly.insertOne(
    {
        keywords: ["Windowfly", "Streifen", "hinterl(ä|ae|a)sst"],
        response: "Überprüfen Sie, ob der Reinigungspad sauber und richtig angebracht ist. Möglicherweise benötigen Sie auch eine andere Art von Reinigungsmitteln."
    }
)

db.windowfly.insertOne(
    {
        keywords: ["Windowfly", "App", "verb(i|u)nde", "nicht"],
        response: "Stellen Sie sicher, dass Ihr Smartphone mit dem gleichen WLAN wie die Windowfly verbunden ist und dass die neueste Version der App installiert ist."
    }
)

db.windowfly.insertOne(
    {
        keywords: ["Windowfly", "stoppt", "pl(ö|oe|o)tzlich"],
        response: "Überprüfen Sie, ob der Akku vollständig geladen ist und ob es auf der Fensterfläche Hindernisse gibt, die das Gerät stoppen könnten."
    }
)

db.windowfly.insertOne(
    {
        keywords: ["Windowfly", "Wasser", "(l(a|ä|ae)|tr(e|i)t).*aus"],
        response: "Stellen Sie sicher, dass der Wasserbehälter richtig eingesetzt ist und keine Risse oder Schäden aufweist. Leeren und reinigen Sie den Behälter regelmäßig."
    }
)

db.cleanbug.insertOne(
    {
        keywords: ["Cleanbug", "start", "nicht"],
        response: "Wenn Ihr Cleanbug nicht startet, überprüfen Sie bitte zuerst, ob er ordnungsgemäß mit Strom versorgt ist. Stellen Sie sicher, dass der Akku vollständig geladen ist und alle Verbindungen fest sitzen. Wenn das Problem weiterhin besteht, könnte es sein, dass ein Sicherheitsschalter nicht richtig aktiviert ist. Überprüfen Sie alle Sicherheitsvorrichtungen und kontaktieren Sie bei Bedarf unseren E-Mail-Support."
    }
)

db.cleanbug.insertOne(
    {
        keywords: ["Cleanbug", "(beweg|fähr|faehr|fahr).*", "nicht"],
        response: "Wenn sich Ihr Cleanbug nicht bewegt, könnte dies auf ein Blockierungsproblem hinweisen. Überprüfen Sie, ob sich etwas im Weg befindet, das die Bewegung des Roboters blockiert. Reinigen Sie die Räder und Sensoren gründlich, da Schmutz und Ablagerungen die Bewegung beeinträchtigen können. Wenn das Problem weiterhin besteht, könnte auch ein Motorschaden vorliegen. In diesem Fall wenden Sie sich bitte an unseren E-Mail-Support."
    }
)

db.cleanbug.insertOne(
    {
        keywords: ["Cleanbug", "warn(leuchte|lampe)", "blink"],
        response: "Eine blinkende Warnleuchte kann verschiedene Probleme signalisieren, darunter einen niedrigen Akkustand, eine blockierte Bürste oder ein Sensorfehler. Überprüfen Sie zuerst den Akkustand und laden Sie den Cleanbug gegebenenfalls auf. Wenn die Bürsten blockiert sind, reinigen Sie diese bitte gründlich. Falls das Problem weiterhin besteht, überprüfen Sie die Bedienungsanleitung auf weitere Hinweise oder wenden Sie sich an unseren E-Mail-Support."
    }
)

db.cleanbug.insertOne(
    {
        keywords: ["Cleanbug", "verl(o|ie)r", "Saug(leistung|kraft)"],
        response: "Wenn Ihr Cleanbug an Saugleistung verliert, könnte dies darauf hinweisen, dass der Staubbehälter voll ist oder dass die Filter verschmutzt sind. Leeren Sie den Staubbehälter regelmäßig aus und reinigen Sie die Filter gemäß den Anweisungen in der Bedienungsanleitung. Stellen Sie sicher, dass auch die Bürsten und Rollen frei von Haaren und Schmutz sind, da dies die Saugleistung beeinträchtigen kann."
    }
)

db.cleanbug.insertOne(
    {
        keywords: ["Cleanbug", "(ein.*schalt|schalt.*ein)", "nicht"],
        response: " Wenn Ihr Cleanbug sich plötzlich nicht mehr einschalten lässt, könnte dies auf einen Defekt im Ein-/Ausschalter oder auf einen Stromversorgungsfehler hinweisen. Überprüfen Sie zuerst den Akkustand und laden Sie den Cleanbug gegebenenfalls auf. Falls das Problem weiterhin besteht, könnte ein Defekt im Ein-/Ausschalter vorliegen. In diesem Fall wenden Sie sich bitte an unseren E-Mail-Support für weitere Unterstützung."
    }
)

db.gardenbeetle.insertOne(
    {
        keywords: ["Gardenbeetle", "Rasen", "ungleichm(ä|ae|a)(ß|ss|s)ig"],
        response: "Stellen Sie sicher, dass die Mähklingen scharf und sauber sind und dass der Gardenbeetle auf die richtige Schnitthöhe eingestellt ist"
    }
)

db.gardenbeetle.insertOne(
    {
        keywords: ["Gardenbeetle", "Akku", "schnell", "leer"],
        response: "Überprüfen Sie, ob der Akku älter ist und ersetzt werden muss, oder ob der Gardenbeetle auf ineffiziente Weise arbeitet. Wenn Sie vermuten, dass es am Akku liegt, wenden Sie sich bitte an unseren E-Mail-Support."
    }
)

db.gardenbeetle.insertOne(
    {
        keywords: ["Gardenbeetle", "Unkraut", "nicht", "entfern"],
        response: "Überprüfen Sie, ob die Unkrauterkennungseinstellungen korrekt sind und ob die Wetterbedingungen die Unkrautentfernung beeinflussen könnte."
    }
)

db.gardenbeetle.insertOne(
    {
        keywords: ["Gardenbeetle", "App", "verb(u|i)nde*", "nicht"],
        response: "Überprüfen Sie, ob Ihr Smartphone und der Gardenbeetle beide mit demselben WLAN verbunden sind und ob die App auf dem neuesten Stand ist."
    }
)

db.gardenbeetle.insertOne(
    {
        keywords: ["Gardenbeetle", "navigiert", "fehlerhaft"],
        response: "Stellen Sie sicher, dass keine Hindernisse die Sensoren blockieren und dass die Software auf dem neuesten Stand ist"
    }
)

db.empty.insertOne(
    {
        keywords: [""],
        response: "Entschuldigen Sie bitte, ich habe Sie nicht verstanden. Könnten Sie das bitte weitere Informationen bereitstellen?\nAm Besten nennen Sie das betroffene Produkt und das konkrete Problem.\n Sollten Sie weitere Unterstützung benötigen, wenden Sie sich bitte an den Telefon- oder E-Mail-Support. Sollten Sie Silber- oder Gold-Status-Kunde sein, können Sie sich natürlich auch an alle weiteren Support-Dienste wenden, die Ihnen zur Verfügung stehen.\n Weitere Informationne finden Sie auf unserer Website."
    }
)

db.empty.insertOne(
    {
        keywords: ["danke"],
        response: "Ich hoffe ich konnte Ihnen helfen. Falls Sie weitere Fragen haben, stehe ich Ihnen gerne zur Verfügung."
    }
)
