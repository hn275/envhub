<script lang="ts">
	import Variable from "./variable";
	import { store } from "../store";
	import Row from "./row.svelte";
	import AddButton from "./new-modal.svelte";
	import { getVariables } from "../services";

	export let repoID: number;
	let rsp: ReturnType<typeof getVariables> = getVariables(repoID);
</script>

<Row>
	<div />
	<h3>Key</h3>
	<h3>Value</h3>
	<h3>Created At</h3>
	<h3>Last Modified</h3>
</Row>

{#await rsp}
	<div class="w-full h-full flex grow justify-center items-center">
		<span class="loading loading-ring loading-lg text-primary" />
	</div>
{:then}
	{#each $store.variables as variable, i (variable.id)}
		<Variable
			{i}
			{...variable}
		/>
	{/each}

	{#if $store.variables.length === 0}
		<div
			class="flex h-full min-h-[400px] w-full flex-col items-center justify-center gap-3"
		>
			<p class="text-light/50">No variables stored</p>
			{#if $store.is_owner}
				<AddButton />
			{/if}
		</div>
	{/if}
{:catch e}
	<div class="w-full h-full flex grow justify-center items-center">
		<div class="text-error">
			<i class="fa-solid fa-circle-exclamation inline" />&nbsp;
			<p class="inline">{e.message}</p>
		</div>
	</div>
{/await}

<style lang="postcss">
	h3 {
		@apply font-semibold text-base-content;
		@apply ml-2;
	}
</style>
