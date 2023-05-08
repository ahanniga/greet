<script>
    import {EventsOn} from "../wailsjs/runtime/runtime.js";

    let readable = 0;
    let writable = 0;
    let colour = 'text-danger';
    let subs = 0;

    const onStatusUpdate = (opts) => {
        readable = opts.readable || 0;
        writable = opts.writable || 0;

        if(readable && writable) {
            colour = 'text-success'
        }
        else if(!readable && !writable) {
            colour = 'text-danger';
        }
        else {
            colour = 'text-warning';
        }

        subs = opts.subs || 0;
    }
    EventsOn("evRelayStatus", onStatusUpdate);

</script>

<style>
</style>

<span class="float-end text-muted">
    <span class="mx-3">|</span>Subs: {subs}
    <span class="mx-3">|</span><i class="bi bi-hdd-network-fill me-2 {colour} "/>R {readable} : W {writable}</span>
