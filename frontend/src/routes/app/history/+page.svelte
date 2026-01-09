<script>
  import ExpensesHistory from '$components/ExpensesHistory.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import { page } from '$app/stores';
  import { getPageCode } from '$lib/utils/pageCodes';

  $: pageCode = getPageCode($page.url.pathname);
</script>

<PageContainer pageCode={pageCode}>
  <ExpensesHistory />
</PageContainer>












