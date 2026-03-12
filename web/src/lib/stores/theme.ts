import { browser } from '$app/environment';
import { writable } from 'svelte/store';

type Theme = 'light' | 'dark';

// 初始化主题
function createThemeStore() {
	// 默认深色主题
	const defaultTheme: Theme = 'dark';
	
	const stored = browser ? (localStorage.getItem('theme') as Theme) : null;
	const prefersDark = browser ? window.matchMedia('(prefers-color-scheme: dark)').matches : true;
	
	const initialTheme = stored || (prefersDark ? 'dark' : 'light');
	
	const { subscribe, set, update } = writable<Theme>(initialTheme);
	
	return {
		subscribe,
		toggle: () => {
			update(current => {
				const newTheme = current === 'dark' ? 'light' : 'dark';
				if (browser) {
					localStorage.setItem('theme', newTheme);
					document.documentElement.classList.toggle('dark', newTheme === 'dark');
				}
				return newTheme;
			});
		},
		setTheme: (theme: Theme) => {
			if (browser) {
				localStorage.setItem('theme', theme);
				document.documentElement.classList.toggle('dark', theme === 'dark');
			}
			set(theme);
		},
		init: () => {
			if (browser) {
				const theme = localStorage.getItem('theme') as Theme || defaultTheme;
				document.documentElement.classList.toggle('dark', theme === 'dark');
				set(theme);
			}
		}
	};
}

export const theme = createThemeStore();