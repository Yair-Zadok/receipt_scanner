<script>
    import { onMount } from 'svelte';
    let index = 0;
 
    let Supplier = "";
    let Account = "";
    let Subtotal = "";
    let Tax = "";
    let Tips = "";
    let Date = "";
    let Encoded_image = "";

    let jsonDataArr; 
    async function load_next_data() {

        console.log(Subtotal);
        if (index != 0) { 
        await fetch('http://localhost:9090/api/post_entry', {
            method: "POST",
            body: JSON.stringify({
                Supplier,
                Account,
                Subtotal,
                Tax,
                Tips,
                Date
            }),
            headers: {
                'Content-Type': 'application/json'
            }
        });
        }

        if (index < jsonDataArr.length) {
        console.log(jsonDataArr);
        let jsonData = jsonDataArr[index];
        console.log(jsonData);

        Account = jsonData.Account;
        Date = jsonData.Date;
        Supplier = jsonData.Supplier;
        Subtotal = jsonData.Subtotal;
        Tax = jsonData.Tax;
        Tips = jsonData.Tips;
        Encoded_image = jsonData.Encoded_image;
        console.log(Encoded_image);
        index += 1;
        }
    }

    async function initialize() {
         const response = await fetch('http://localhost:9090/api/data');
         jsonDataArr = await response.json();
    }

    onMount(async function() {
        await initialize();
        load_next_data();
    });


</script>

<div style="display: flex;">
    <div style="flex: 1;">
        <form on:submit|preventDefault={load_next_data}>
            <label>Supplier:<br>
                <input type="text" bind:value={Supplier}> 
            </label>
            <label>Account:<br>
                <input type="text" bind:value={Account}> 
            </label>
            <label>Subtotal:<br>
                <input type="text" bind:value={Subtotal}> 
            </label>
            <label>Tax:<br>
                <input type="text" bind:value={Tax}> 
            </label>
            <label>Tips:<br>
                <input type="text" bind:value={Tips}> 
            </label>
            <label>Date:<br>
                <input type="text" bind:value={Date}> 
            </label>

            <br><button type="submit">Submit</button>
            <button on:click={initialize} type="button">Update</button>
        </form>
    </div>

    <div style="flex: 1;">
        <img src={"data:image/png;base64, " + Encoded_image} alt="Base64 Image" style="max-width: 500px; max-height: 700px;">
    </div>
</div>

