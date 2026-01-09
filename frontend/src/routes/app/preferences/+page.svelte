<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { theme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';
  import { requireAuthWithSleep } from '$lib/utils/authSleep';

  const availableThemes = ['cupcake', 'night','valentine','dracula'];
  let selectedTheme = 'cupcake';
  
  $: pageCode = getPageCode($page.url.pathname);

  onMount(async () => {
    const ok = await requireAuthWithSleep();
    if (!ok) return;
    selectedTheme = $theme || 'cupcake';
  });

  $: if ($theme) {
    selectedTheme = $theme;
  }

  function handleThemeChange() {
    theme.setTheme(selectedTheme);
  }
</script>

<div class="max-w-4xl mx-auto space-y-4">
  <PageHeader
    title="Preferences"
    subtitle="Manage how the app looks and behaves for your account."
    pageCode={pageCode}
  />
  <div class="card bg-base-100 shadow-xl border-1">
    <div class="card-body space-y-4">
      <!-- Navigation links in two columns -->
      <div class="grid gap-3 md:grid-cols-2">
        <a href="/app/preferences/categories"
           class="btn btn-soft justify-between w-full normal-case h-full min-h-[120px] text-base">
          <span class="text-left">
            <span class="block font-semibold">Categories</span>
            <span class="block text-xs text-base-content/70">
              Manage and reorder your income & expense categories.
            </span>
          </span>
          <span class="text-lg">›</span>
        </a>

        <a href="/app/preferences/currency"
           class="btn btn-soft justify-between w-full normal-case h-full min-h-[120px] text-base">
          <span class="text-left">
            <span class="block font-semibold">Currency & Quick Amount</span>
            <span class="block text-xs text-base-content/70">
              Configure display currency and quick amount shortcuts.
            </span>
          </span>
          <span class="text-lg">›</span>
        </a>

        <a href="/app/preferences/change-password"
           class="btn btn-soft justify-between w-full normal-case h-full min-h-[120px] text-base">
          <span class="text-left">
            <span class="block font-semibold">Change Password</span>
            <span class="block text-xs text-base-content/70">
              Update your account password securely.
            </span>
          </span>
          <span class="text-lg">›</span>
        </a>
      </div>

      <!-- Appearance -->
      <div class="divider my-2"></div>

      <div class="space-y-2 max-w-md">
        <h3 class="font-semibold text-base">Appearance</h3>
        <p class="text-xs text-base-content/70">
          Choose a DaisyUI theme for the app.
        </p>
        <fieldset class="fieldset">
          <legend class="fieldset-legend">Theme</legend>
          <select
            class="select select-bordered w-full border-2"
            bind:value={selectedTheme}
            on:change={handleThemeChange}
          >
            {#each availableThemes as t}
              <option value={t}>{t}</option>
            {/each}
          </select>
        </fieldset>
      </div>
    </div>
  </div>
</div>