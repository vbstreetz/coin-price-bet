import './main.css';
import 'nprogress/nprogress.css';
import 'bulma/css/bulma.css';
import 'tailwindcss/dist/base.css';
import 'tailwindcss/dist/components.css';
import 'tailwindcss/dist/utilities.css';

import HMR from '@sveltech/routify/hmr';
import App from './App.svelte';

const app = HMR(App, { target: document.body }, 'routify-app');

export default app;
