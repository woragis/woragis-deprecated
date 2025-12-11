export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
	id: string;
	message: string;
	type: ToastType;
	duration?: number;
}

class ToastManager {
	private toasts: Toast[] = [];
	private listeners: Set<(toasts: Toast[]) => void> = new Set();

	subscribe(callback: (toasts: Toast[]) => void) {
		this.listeners.add(callback);
		callback(this.toasts);
		return () => {
			this.listeners.delete(callback);
		};
	}

	private notify() {
		this.listeners.forEach(callback => callback(this.toasts));
	}

	show(message: string, type: ToastType = 'success', duration = 3000) {
		const id = Math.random().toString(36).substr(2, 9);
		const toast: Toast = { id, message, type, duration };
		this.toasts = [...this.toasts, toast];
		this.notify();
		return id;
	}

	remove(id: string) {
		this.toasts = this.toasts.filter(t => t.id !== id);
		this.notify();
	}

	clear() {
		this.toasts = [];
		this.notify();
	}

	getToasts(): Toast[] {
		return this.toasts;
	}
}

export const toastManager = new ToastManager();

export function showToast(message: string, type: ToastType = 'success', duration = 3000) {
	return toastManager.show(message, type, duration);
}
