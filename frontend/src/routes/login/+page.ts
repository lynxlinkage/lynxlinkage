// Login is dynamic, authenticated, never prerendered. SSR is also disabled
// so the page renders in the browser only — that keeps cookies, fetch, and
// form state firmly on the client.
export const prerender = false;
export const ssr = false;
