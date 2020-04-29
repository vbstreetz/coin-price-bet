import { API_HOST } from '../config';

export async function query(endpoint) {
  const res = await fetch(buildUrl(endpoint));
  const { result } = await res.json();
  return result;
}

export async function mutate(endpoint, data) {
  const res = await fetch(buildUrl(endpoint), { method: 'POST', body: data });
  return res.json();
}

function buildUrl(path) {
  return `${API_HOST}:1317/${path}`;
}
