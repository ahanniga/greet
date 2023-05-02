<script>
    import {onMount} from "svelte";

    let text = "";

    onMount(() => {
        document.getElementById("statusMessage").hidden = true;
    })

    const onStatusText = (txt) => {
        let d = document.getElementById("statusMessage");
        d.hidden = false;
        text = txt;
        setTimeout(() => {
            text = "";
            d.hidden = true;
        }, 5000);
    }
    window.runtime.EventsOn('evStatusText', onStatusText);

</script>
<style></style>

<div id="statusMessage" class="col-12 py-2">
    <div class="alert alert-secondary" role="alert">{text}</div>
</div>