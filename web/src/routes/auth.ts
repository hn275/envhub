import { makeUrl } from "@lib/url";
import type { User } from "@lib/auth";

export async function signIn(code: string): Promise<User> {
	const res = await fetch(makeUrl("/auth/github"), {
		method: "POST",
		headers: {
			"content-type": "application/json",
			accept: "application/json",
		},
		body: JSON.stringify({ code }),
	});
	const payload = await res.json();
	if (res.status !== 200) throw new Error(payload["message"]);
	return payload as User;
}
