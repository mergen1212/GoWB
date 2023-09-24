export async function load({params}){
    async function GetPost(){
        let url = "http://localhost:8080/Getitemstotrek/WBILMTESTTRACK/"+params.id
        const res = await fetch(url)
        let data = res.json()
        return data
    }
    return{
        Items: await GetPost()
    }
}