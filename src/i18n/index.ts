export const languages = ["en", "de", "es"] as const;
export type Lang = (typeof languages)[number];

export const defaultLang: Lang = "en";

export const labels: Record<Lang, string> = {
	en: "EN",
	de: "DE",
	es: "ES",
};

export type Content = typeof contentEn;

const contentEn = {
	nav: {
		features: "What",
		manifesto: "Manifesto",
	},
	hero: {
		line1: "Real Music.",
		line2: "Real People.",
		notPlatform: "Not a streaming platform.",
		notTalent: "Not a talent show.",
		workshop: "A workshop.",
		description:
			"Muziboo is a space for hobby musicians, bedroom producers, and anyone who makes music with their own hands.",
		share: "Share your sound with people who care about craft — not metrics.",
	},
	features: {
		title: "A workshop, not a stage.",
		items: [
			{
				title: "Upload & Share",
				body: "Demos. Fragments. Late-night recordings. Unfinished songs. Share what you are working on with people who actually listen.",
			},
			{
				title: "Human-First Culture",
				body: "Every track here begins with a person. Not because we reject technology — but because we believe music starts with intention.",
			},
			{
				title: "Work in Progress",
				body: "Music does not arrive finished. We are building tools that let you show how a song evolves — from first take to final mix.",
			},
			{
				title: "Discover Real",
				body: "Focused discovery. Transparent timelines. No ranking games. No invisible systems deciding who gets heard.",
			},
			{
				title: "Feedback Threads",
				body: "Real conversations about the work. Constructive, public, and grounded in listening — not likes and follows.",
			},
			{
				title: "Your Space",
				body: "Build your catalog. Tell your story. A profile that reflects your journey — not an engagement strategy.",
			},
		],
	},
	manifesto: {
		title: "Muziboo Manifesto — 2026",
		paragraphs: [
			"The internet was not built for perfection. It was built for participation.",
			"Before streams were counted and feeds were ranked, people shared what they were working on. Demos. Fragments. Late-night recordings. Unfinished songs.",
			"Music was conversation.",
			"Muziboo is a home for that conversation.",
			"We believe music is not content. It is craft.",
			"It is the sound of someone trying. Learning. Improving. Connecting.",
			"We are not here for virality. We are not here for metrics. We are not here for shortcuts.",
			"We are here for musicians who make music because they must.",
			"Bedroom producers. Garage bands. Choir singers. Guitarists who started last month. Drummers who have been playing for thirty years.",
			"Every note here was played by someone who cared enough to try.",
			"No ranking games. No artificial amplification. No invisible systems deciding who gets heard.",
			"Chronological. Transparent. Human.",
			"Muziboo is not a stage.",
			"It is a workshop. A rehearsal room. A community hall. A place where imperfect things are welcome.",
			"Because music does not begin polished.",
			"It begins with people.",
			"And people are enough.",
		],
	},
	signup: {
		title: "Join the Waitlist",
		name: "Name",
		email: "Email Address",
		submit: "Submit",
		button: "Be one of the first",
		message: "What do you play? Tell us more... (Optional)",
		successTitle: "Thank you!",
		successText: "You've been added to the list. We will be in touch soon.",
	},
	skipToContent: "Skip to main content",
	footer: {
		brand: "muziboo",
		est: "EST. 2007",
		reborn: "REBORN 2026",
		facebook: "FACEBOOK",
		privacy: "PRIVACY",
		terms: "TERMS",
		copyright: "© 2026 MUZIBOO",
	},
};

const contentDe: Content = {
	nav: {
		features: "Was",
		manifesto: "Manifest",
	},
	hero: {
		line1: "Real Music.",
		line2: "Real People.",
		notPlatform: "Kein Wettbewerb um Aufmerksamkeit.",
		notTalent: "Ohne Algorithmen, die entscheiden.",
		workshop: "Ein Ort für Musiker.",
		description:
			"Muziboo ist für Hobbymusiker, Homestudio-Producer und alle, die ihre Musik selbst machen.",
		share:
			"Teile deinen Sound mit anderen Menschen, aus Liebe zur Musik — nicht für Algorithmen.",
	},
	features: {
		title: "Ein Proberaum, keine Bühne.",
		items: [
			{
				title: "Hochladen & Teilen",
				body: "Demos. Fragmente. Aufnahmen mitten in der Nacht. Unfertige Songs. Teile, woran du arbeitest, mit Menschen, die wirklich zuhören.",
			},
			{
				title: "Mensch zuerst",
				body: "Jeder Track entspringt einem Menschen. Wir lieben Software – aber die Intention, das Gefühl, kommt vom Musiker. Darum steht bei uns der Mensch im Mittelpunkt.",
			},
			{
				title: "Work in Progress",
				body: "Musik entsteht Schritt für Schritt. Auf Muziboo kannst du zeigen, wie ein Song wächst – vom ersten Take bis zum finalen Mix. Teile jede Phase und tausche dich mit anderen darüber aus.",
			},
			{
				title: "Echtes Entdecken",
				body: "Musik erscheint bei uns so, wie sie hochgeladen wird – chronologisch und transparent. Es gibt keine Rankings und keine versteckten Algorithmen, die entscheiden, wer gehört wird.",
			},
			{
				title: "Dialog",
				body: "Austausch über Musik. Offen, konstruktiv und geprägt vom Zuhören – nicht von Likes oder Followern.",
			},
			{
				title: "Dein Raum",
				body: "Lass deine Musik sprechen. Sammle deine Stücke und erzähle deinen Weg. Dein Profil zeigt, wer du bist – nicht wie gut du dich vermarktest",
			},
		],
	},
	manifesto: {
		title: "Muziboo Manifest — 2026",
		paragraphs: [
			"Das Internet wurde nicht für Perfektion gebaut, sondern fürs Mitmachen.",
			"Bevor Streams gezählt und Feeds sortiert wurden, teilten Menschen, woran sie arbeiteten.",
			"Demos. Skizzen. Aufnahmen mitten in der Nacht.",
			"Unfertige Songs.",
			"Musik war Austausch.",
			"Muziboo ist ein neuer alter Ort für genau diesen Austausch.",
			"Musik ist mehr als ein Beitrag im Feed, sie ist persönlich.",
			"Musik ist der Klang von jemandem, der ausprobiert, lernt, weitergeht.",
			"Bei uns geht es nicht um Reichweite, nicht um Zahlen.",
			"Hier geht es um Musik.",
			"Für Menschen, die Musik machen, weil sie es müssen.",
			"Weil sie nicht anders können.",
			"Homestudio-Producer. Bands im Proberaum. Chorsänger. Gitarristen, die gerade erst anfangen. Schlagzeuger mit dreißig Jahren Erfahrung.",
			"Jede Note kommt von jemandem, dem sie etwas bedeutet.",
			"Keine Rankings. Keine künstliche Verstärkung. Keine Systeme, die entscheiden, wer sichtbar ist.",
			"Einfach Musik.",
			"In der Reihenfolge, in der sie entsteht.",
			"Muziboo ist kein Wettbewerb.",
			"Es ist ein Raum. Für Ideen. Für Entwicklung. Für Unfertiges.",
			"Denn Musik beginnt nicht poliert.",
			"Sie beginnt mit Menschen.",
			"Und das reicht.",
		],
	},
	signup: {
		title: "Auf die Warteliste",
		name: "Name",
		email: "E-Mail-Adresse",
		submit: "Absenden",
		button: "Als Erste dabei sein",
		message: "Was spielst du? Erzähl uns mehr... (Optional)",
		successTitle: "Danke!",
		successText: "Du stehst auf der Liste. Wir melden uns in Kürze.",
	},
	skipToContent: "Zum Hauptinhalt",
	footer: {
		brand: "muziboo",
		est: "EST. 2007",
		reborn: "REBORN 2026",
		facebook: "FACEBOOK",
		privacy: "DATENSCHUTZ",
		terms: "NUTZUNGSBEDINGUNGEN",
		copyright: "© 2026 MUZIBOO",
	},
};

const contentEs: Content = {
	nav: {
		features: "Que",
		manifesto: "Manifiesto",
	},
	hero: {
		line1: "Música real.",
		line2: "Gente real.",
		notPlatform: "No es una plataforma de streaming.",
		notTalent: "No es un concurso de talentos.",
		workshop: "Un taller.",
		description:
			"Muziboo es un espacio para músicos aficionados, productores de habitación y cualquiera que haga música con sus propias manos.",
		share: "Comparte tu sonido con quien valora el oficio — no las métricas.",
	},
	features: {
		title: "Un taller, no un escenario.",
		items: [
			{
				title: "Sube y comparte",
				body: "Demos. Fragmentos. Grabaciones de madrugada. Canciones sin terminar. Comparte en qué trabajas con quien de verdad escucha.",
			},
			{
				title: "Cultura centrada en las personas",
				body: "Cada tema aquí empieza con una persona. No porque rechacemos la tecnología — sino porque creemos que la música empieza con la intención.",
			},
			{
				title: "Trabajo en progreso",
				body: "La música se construye poco a poco. En Muziboo puedes compartir cómo evoluciona una canción — desde el primer ensayo hasta la mezcla final — y hablar del proceso con otros músicos.",
			},
			{
				title: "Descubre de verdad",
				body: "Aquí nadie decide qué merece ser escuchado. La música está una junto a otra, como fue creada.",
			},
			{
				title: "Diálogo",
				body: "Conversaciones reales sobre el trabajo. Constructivas, públicas y basadas en escuchar — no en likes ni seguidores.",
			},
			{
				title: "Tu espacio",
				body: "Deja que tu música hable. Reúne tus canciones y cuenta tu camino. Tu perfil muestra quién eres — no lo bien que sabes promocionarte.",
			},
		],
	},
	manifesto: {
		title: "Manifiesto Muziboo — 2026",
		paragraphs: [
			"Internet no nació para la perfección.",
			"Nació para participar.",
			"",
			"Antes de que se contaran los streams y se ordenaran los feeds, la gente compartía en qué estaba trabajando.",
			"",
			"Demos. Bocetos. Grabaciones a medianoche.",
			"",
			"Canciones sin terminar.",
			"",
			"La música era conversación.",
			"",
			"Muziboo es un lugar — nuevo y antiguo a la vez — para recuperar esa conversación.",
			"",
			"La música es más que una publicación en un feed.",
			"Es algo personal.",
			"",
			"Es el sonido de alguien que prueba, aprende y sigue adelante.",
			"",
			"Aquí no se trata de alcance.",
			"Ni de cifras.",
			"",
			"Aquí se trata de música.",
			"",
			"Para quienes hacen música porque lo necesitan.",
			"Porque no pueden no hacerlo.",
			"",
			"Productores desde casa. Bandas en el local de ensayo. Coros. Guitarristas que acaban de empezar. Bateristas con treinta años de experiencia.",
			"Cada nota viene de alguien a quien le importa.",
			"Sin rankings.",
			"Sin impulso artificial.",
			"Sin sistemas que decidan quién merece ser escuchado.",
			"",
			"Solo música.",
			"",
			"En el orden en que va naciendo.",
			"",
			"Muziboo no es una competencia.",
			"",
			"Es un espacio.",
			"",
			"Para ideas.",
			"Para crecer.",
			"Para lo inacabado.",
			"",
			"Porque la música no empieza pulida.",
			"",
			"Empieza con personas.",
			"",
			"Y eso basta.",
		],
	},
	signup: {
		title: "Únete a la lista de espera",
		name: "Nombre",
		email: "Correo electrónico",
		submit: "Enviar",
		button: "Sé de los primeros",
		message: "¿Qué tocas? Cuéntanos más... (Opcional)",
		successTitle: "¡Gracias!",
		successText: "Te hemos añadido a la lista. Estaremos en contacto pronto.",
	},
	skipToContent: "Saltar al contenido",
	footer: {
		brand: "muziboo",
		est: "EST. 2007",
		reborn: "REBORN 2026",
		facebook: "FACEBOOK",
		privacy: "PRIVACIDAD",
		terms: "TÉRMINOS",
		copyright: "© 2026 MUZIBOO",
	},
};

export const content: Record<Lang, Content> = {
	en: contentEn,
	de: contentDe,
	es: contentEs,
};
