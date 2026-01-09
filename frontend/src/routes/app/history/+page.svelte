<script>
  import ExpensesHistory from '$components/ExpensesHistory.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import CurrentWorkspaceBadge from '$lib/components/CurrentWorkspaceBadge.svelte';
  import { page } from '$app/stores';
  import { getPageCode } from '$lib/utils/pageCodes';

  $: pageCode = getPageCode($page.url.pathname);
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title="History"
    subtitle="Review your past income and expenses."
    pageCode={pageCode}
  />
  <CurrentWorkspaceBadge />
  <ExpensesHistory />
</PageContainer>












