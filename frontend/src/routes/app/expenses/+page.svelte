<script>
  import { page } from '$app/stores';
  import InputExpenses from '$components/InputExpenses.svelte';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';

  $: pageCode = getPageCode($page.url.pathname);
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title="Add Expense"
    subtitle="Record your expenses and track your spending."
    pageCode={pageCode}
  />
  <div class="card bg-base-100 shadow-xl border border-base-300">
    <InputExpenses/>
  </div>
</PageContainer>
