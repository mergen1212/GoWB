<script>
    import TitleBlock from "./[id]/TitleBlock.svelte";

    let searchQuery = '';
    const handleSearch = async () => {
        // Выполняем поиск или другую логику с использованием searchQuery
        console.log('Поиск выполнен с запросом:', searchQuery);

        // Выполняем GET-запрос на сервер с использованием searchQuery
        async function GetPosts(){
            let url = `http://localhost:8080/Getitemstotrek/${searchQuery}`
            const res = await fetch(url)
            let data = res.json()
            console.log(data)
            return data
        }
        return{
            Items: await GetPosts()
        }
    };
    const handleSearchOrder = async () => {
        // Выполняем поиск или другую логику с использованием searchQuery
        console.log('Поиск выполнен с запросом:', searchQuery);

        // Выполняем GET-запрос на сервер с использованием searchQuery
        async function GetPosts(){
            let url = `http://localhost:8080/GetOrderUid/${searchQuery}`
            const res = await fetch(url)
            let data = res.json()
            console.log(data)
            return data
        }
        return{
            Order: await GetPosts()
        }
    };
</script>

<form on:submit={handleSearch}>
    <label>
        Поиск items по track_number:
        <input type="text" bind:value={searchQuery} />
    </label>
    <button type="submit">Искать</button>
</form>


<form on:submit={handleSearchOrder}>
    <label>
        Поиск order по order_uid:
        <input type="text" bind:value={searchQuery} />
    </label>
    <button type="submit">Искать</button>
</form>

