<script>
    import TitleBlock from "./TitleBlock.svelte";

    export let data
</script>

<h1>post</h1>
<TitleBlock>
    {#each data.Items as item}
        <p>Item id: {item.id}</p>
        <p>chrt_id: {item.chrt_id}</p>
        <p>price: {item.price}</p>
        <p>rid: {item.rid}</p>
        <p>name: {item.name}</p>
        <p>sale: {item.sale}</p>
        <p>size: {item.size} </p>
        <p>total_price: {item.total_price}</p>
        <p>nm_id: {item.nm_id}</p>
        <p>brand: {item.brand}</p>
        <p>track_number:{item.track_number}</p>
        <p>createdAt: {item.createdAt}</p>
        <p></p>
    {/each}
</TitleBlock>