type ConfirmOptions = {
	title: string;
	description?: string;
	action?: string;
	destructive?: boolean;
};

type Pending = ConfirmOptions & { resolve: (ok: boolean) => void };

class ConfirmState {
	pending = $state<Pending | null>(null);

	/** Shows the confirmation dialog and resolves with the user's answer. */
	ask(o: ConfirmOptions): Promise<boolean> {
		this.pending?.resolve(false);
		return new Promise((resolve) => (this.pending = { ...o, resolve }));
	}

	/** Closes the dialog and resolves the pending ask(). */
	answer(ok: boolean) {
		this.pending?.resolve(ok);
		this.pending = null;
	}
}

export const confirmState = new ConfirmState();

/** Opens the confirmation dialog and resolves to true when the user confirms. */
export const ask = (o: ConfirmOptions) => confirmState.ask(o);
