import { writable } from "svelte/store";

export type NewVariable = {
	key: string;
	value: string;
};

export interface Variable extends NewVariable {
	id: string;
	created_at: Date;
	updated_at: Date;
}

export type RepositoryEnv = {
	write_access: boolean;
	variables: Array<Variable>;
};

export const store = writable<RepositoryEnv>({
	write_access: false,
	variables: [],
});
