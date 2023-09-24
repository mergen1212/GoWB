export async function load(){
    async function GetPosts(){
        let url = "http://localhost:8080/Getitemstotrek/WBILMTESTTRACK"
        const res = await fetch(url)
        let data = res.json()
        return data
    }
    return{
        Items: await GetPosts()
    }
}