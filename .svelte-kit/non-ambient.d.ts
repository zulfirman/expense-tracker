
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/budget" | "/categories" | "/expenses" | "/history" | "/income" | "/login" | "/profile" | "/signup";
		RouteParams(): {
			
		};
		LayoutParams(): {
			"/": Record<string, never>;
			"/budget": Record<string, never>;
			"/categories": Record<string, never>;
			"/expenses": Record<string, never>;
			"/history": Record<string, never>;
			"/income": Record<string, never>;
			"/login": Record<string, never>;
			"/profile": Record<string, never>;
			"/signup": Record<string, never>
		};
		Pathname(): "/" | "/budget" | "/budget/" | "/categories" | "/categories/" | "/expenses" | "/expenses/" | "/history" | "/history/" | "/income" | "/income/" | "/login" | "/login/" | "/profile" | "/profile/" | "/signup" | "/signup/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): string & {};
	}
}