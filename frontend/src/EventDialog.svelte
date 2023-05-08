<script>
    /**
     *  Display a note in a dialog.
     *  TODO: Remove. Display inline with the parent event
     */

    import {GetMyPubkey, GetTextNotesByEventIds, Nip19Decode} from "../wailsjs/go/main/App.js";
    import {EventsOn} from "../wailsjs/runtime/runtime.js";
    import EventPost from "./EventPost.svelte";

    let myPk;
    let event = false;

    const onEventDialog = (noteRef) => {
        document.getElementById('launchEventDialog').click();
        event = false;
        GetMyPubkey().then((pk) => {
            myPk = pk;
            if(noteRef.startsWith("note" || noteRef.startsWith("nevent"))) {
                Nip19Decode(noteRef).then((hexId)=>{
                   GetTextNotesByEventIds([hexId]).then((events)=>{
                       if(events.length > 0) {
                           event = events[0];
                       }
                   });
                }).catch((err)=>{
                    console.log("Error:" + err);
                });
            } else {
                GetTextNotesByEventIds([noteRef]).then((events)=>{
                    if(events.length > 0) {
                        event = events[0];
                    }
                });
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
                <h1 class="modal-title fs-5" id="staticBackdropLabel">Event {event.id}</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                {#if event}
                    <EventPost {event} {myPk}/>
                {/if}
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>
