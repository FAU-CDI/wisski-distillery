import Call, { Hooks } from './client'
import { Provision, Rebuild } from './client/calls';

const eConsole = new console.Console(process.stderr, process.stderr);

// read API KEY 
const API_KEY = process.env.API_KEY;
if (!API_KEY) {
    eConsole.error('API_KEY not speciied')
}

// READ ARGUMENTS
if (process.argv.length < 4) {
    eConsole.error('Usage: API_KEY=$API_KEY <script> $REMOTE $SLUG');
    process.exit(1);
}
const REMOTE = process.argv[2];
const SLUG = process.argv[3];

// do the call!
const result = Call(
    {
        url: REMOTE,
        token: API_KEY,
    },
    Rebuild(SLUG, {
        PHP: "Default (8.1)",
        OpCacheDevelopment: false,
        ContentSecurityPolicy: "",
    }),
    {
        beforeCall: eConsole.log.bind(eConsole, 'beforeCall'),
        afterCall: eConsole.log.bind(eConsole, 'afterCall'),
        onError: eConsole.error.bind(eConsole, 'onError'),
        onLogLine: (_, line) => process.stdout.write(line)
    },
);

result.then((x) => eConsole.log(x)).catch((x) => eConsole.error(x));