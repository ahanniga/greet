<script>
    let title = "";
    let message = "";
    let cancelable = false;
    let okButton = "OK";
    let iconClass = "bi-info-circle"
    const onMessageDialog = (opts) => {
        title = opts.title || "Message";
        message = opts.message || "";
        cancelable = opts.cancelable || false;
        okButton = opts.okButton || "OK";
        iconClass = opts.iconClass || "bi-info-circle";
        document.getElementById('launchMessageDialog').click();
    }
    window.runtime.EventsOn('evMessageDialog', onMessageDialog);

</script>
<style></style>

<a class="visually-hidden" id="launchMessageDialog" data-bs-toggle="modal" data-bs-target="#messageDialog"></a>
<div class="modal fade" id="messageDialog" data-bs-backdrop="static" >
    <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="messageDialogTitle"><i class="bi {iconClass} me-3"/>{title}</h5>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <p>{@html message}</p>
            </div>
            <div class="modal-footer">
                {#if cancelable}
                    <button id="boostClose" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                {/if}
                <button type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">{okButton}</button>
            </div>
        </div>
    </div>
</div>