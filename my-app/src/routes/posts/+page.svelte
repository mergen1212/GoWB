<script>
    import OrderForm from "./OrderForm.svelte";

    export let data;
</script>

<h1>Posts</h1>

<OrderForm/>

{#each data.Items as Item}
  <p>
      <a class="link" href="/posts/{Item.track_number}">
          {Item.track_number}
      </a>
  </p>
{/each}
<style>
    .link{
        -webkit-text-fill-color: white;
    }

</style>