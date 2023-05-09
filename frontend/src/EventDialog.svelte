<script>
    /**
     *  Display a note in a dialog.
     *  TODO: Remove. Display inline with the parent event
     */

    import {GetMyPubkey, GetTextNotesByEventIds, Nip19Decode} from "../wailsjs/go/main/App.js";
    import {EventsOn, LogInfo} from "../wailsjs/runtime/runtime.js";
    import EventPost from "./EventPost.svelte";
    import loadingGif from "./assets/images/loading.gif";

    let myPk;
    let eventId = "";
    let promise = false;

    const onEventDialog = (noteRef) => {
        LogInfo("onEventDialog...");
        eventId = noteRef;
        promise = false;
        document.getElementById('launchEventDialog').click();
        GetMyPubkey().then((pk) => {
            myPk = pk;
            if(noteRef.startsWith("note" || noteRef.startsWith("nevent"))) {
                Nip19Decode(noteRef).then((parts)=>{
                   promise = GetTextNotesByEventIds([parts[1]]);
                }).catch((err)=>{
                    console.log("Error:" + err);
                });
            } else {
                promise = GetTextNotesByEventIds([noteRef]);
            }
        });
    }
    EventsOn('evEventDialog', onEventDialog);

</script>
<style></style>

<a class="visually-hidden" id="launchEventDialog" data-bs-toggle="modal" data-bs-target="#eventDialog"></a>
<div class="modal " id="eventDialog" tabindex="-1" data-bs-backdrop="static" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel">Event {eventId}</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                {#if promise}
                    {#await promise}
                        <img src="{loadingGif}" width="18" height="18">
                    {:then events}
                        {#each events as event}
                        <EventPost {event} {myPk}/>
                        {/each}
                    {/await}
                {/if}
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>
