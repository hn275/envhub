<script lang="ts">
	import { createEventDispatcher, onMount } from "svelte";
	import { handleEdit } from "../../services";
	import { store } from "../../store";

	export let saveAble: boolean;
	export let key: string;
	export let value: string;
	export let newKey: string;
	export let newValue: string;
	export let id: string;
	export let _class: string;

	let disabled: boolean;
	$: disabled = !saveAble;

	$: state = $store;

	let modal: HTMLDialogElement;
	onMount(() => {
		modal = document.getElementById("edit-variable") as HTMLDialogElement;
	});

	const dispatch = createEventDispatcher<{ success: void }>();

	let err: string | undefined;
	let loading: boolean = false;
	async function handleSubmit() {
		try {
			if (!state.repoID) throw new Error("Respository ID not found.");
			loading = true;
			await handleEdit(state.repoID, id, newKey, newValue);
			modal?.close();
			dispatch("success");
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	}
</script>

<button
	class={_class}
	on:click={() => modal?.showModal()}
	{disabled}
>
	<i class="fa-solid fa-check fa-sm" />
</button>

<dialog
	id="edit-variable"
	class="modal"
>
	<form
		method="dialog"
		class="modal-box"
	>
		<h3 class="mb-5 text-lg font-bold">Edit variable</h3>
		<div class="h-44">
			<p class="mb-2">Apply the fowllowing changes:</p>

			<div class="mb-3 flex items-center gap-3">
				<h6 class="w-[4ch] font-semibold">Key</h6>
				<p class="bg-neutral flex items-center gap-3 rounded-md p-3">
					{#if key !== newKey}
						<span>{key}</span>
						<i class="fa-solid fa-arrow-right" />
						<span>{newKey}</span>
					{:else}
						<span>{key}</span>
					{/if}
				</p>
			</div>

			<div class="flex items-center gap-3">
				<h6 class="w-[4ch] font-semibold">Value</h6>
				<p class="bg-neutral flex items-center gap-3 rounded-md p-3">
					{#if value !== newValue}
						<span>{value}</span>
						<i class="fa-solid fa-arrow-right" />
						<span>{newValue}</span>
					{:else}
						<span>{value}</span>
					{/if}
				</p>
			</div>

			{#if err}
				<p class="text-error mt-2">{err}</p>
			{/if}
		</div>

		<div class="mt-3 flex justify-end gap-3">
			<button
				class="btn btn-ghost"
				type="button"
				on:click={() => modal?.close()}
			>
				Cancel
			</button>

			<button
				class="btn btn-primary w-20"
				on:click|preventDefault={handleSubmit}
				disabled={loading}
			>
				{#if loading}
					<span class="loading loading-xs"></span>
				{:else}
					Submit
				{/if}
			</button>
		</div>
	</form>
	<form
		id="edit-var-close"
		method="dialog"
		class="modal-backdrop"
	>
		<button id="edit-var-close-btn" />
	</form>
</dialog>
