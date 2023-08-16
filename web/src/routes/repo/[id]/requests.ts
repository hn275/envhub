import { apiFetch } from "@lib/requests";
import { makeUrl } from "@lib/url";
import {
	store,
	type RepositoryEnv,
	type NewVariable,
	type Variable,
} from "./store";

export async function getVariables(repoID: number): Promise<void> {
	const url = makeUrl(`/repos/${repoID}`);
	const headers = new Headers({
		Accept: "application/json",
	});
	const rsp = await apiFetch(url, { method: "GET", headers });
	const payload: EnvHub.Response<RepositoryEnv> = await rsp.json();
	if (rsp.status !== 200) {
		throw new Error((payload as EnvHub.Error).message);
	}
	store.set(payload as RepositoryEnv);
}

export async function writeNewVariable(
	repoID: number,
	v: NewVariable,
): Promise<void> {
	// mocking a request
	return new Promise((res) => {
		console.log(repoID, v);
		setTimeout(() => {
			const env: Variable = {
				id: String(Math.random() * 100),
				key: "foo",
				value: "bar",
				created_at: new Date(),
				updated_at: new Date(),
			};
			store.update((v) => ({ ...v, variables: [env, ...v.variables] }));
			res();
		}, 2000);
	});
	/*
  const url: string = makeUrl(`/repos/${repoID}`);
  const headers: Headers = new Headers({
    Accept: "application/json",
    "Content-type": "application/json",
  });
  const body: BodyInit = JSON.stringify(v);
  const rsp = await apiFetch(url, { method: "POST", headers, body });
  if (rsp.status !== 201) {
    const payload: EnvHub.Error = await rsp.json();
    throw new Error(payload.message);
  }
  const payload: Variable = await rsp.json();
  store.update((v) => ({ ...v, variables: [payload, ...v.variables] }));
  */
}