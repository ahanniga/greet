<script>
    let event = {};

    const onEventInfo = (ev) => {
        event = ev;
    }
    window.runtime.EventsOn('evEventInfo', onEventInfo);

</script>
<style></style>

<div class="modal " id="eventInfo" tabindex="-1" aria-labelledby="staticBackdropLabel" data-bs-backdrop="static" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel"><i class="bi-info-circle text-muted me-2"></i> Event Info</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <pre class="text-success-emphasis">{JSON.stringify(event, undefined, 2)}</pre>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>