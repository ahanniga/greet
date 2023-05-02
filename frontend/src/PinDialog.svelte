<script>
    import {LoginWithPin} from "../wailsjs/go/main/App.js";

    const onPinDialog = () => {
        setTimeout(() => {
            document.getElementById("pin").value = "";
            document.getElementById('pin').focus();
        }, 500);
    }
    window.runtime.EventsOn('evPinDialog', onPinDialog);

    const showError = (msg) => {
        let d = document.getElementById("pinDialogErrorMessage");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
        }, 5000);
    }

    const setPin = () => {
        let pin = document.getElementById("pin").value;
        LoginWithPin(pin).then((e) => {
            console.log("LoginWithPin returns " + e);
            document.getElementById('pinDialogClose').click();
        }).catch((e) => {
            console.error(e);
            showError(e);
        });
    }

    const checkForEnter = (ev) =>{
        if (ev.charCode === 13) {
            setPin();
        }
    }

    const resetAccount = () => {
        window.runtime.EventsEmit("evLoginDialog");
    }

</script>
<style></style>

<div class="modal" id="pinDialog" tabindex="-1" aria-hidden="true" data-bs-backdrop="static" data-bs-keyboard="false">
    <div class="modal-dialog modal-dialog-centered " >
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" ><i class="bi-box-arrow-in-right me-3"></i>Welcome to Greet!</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <div class="row g-3 align-items-center">
                    <div class="col-auto">
                        <label for="pin" class="col-form-label">PIN</label>
                    </div>
                    <div class="col-auto">
                        <input type="password" id="pin" on:keypress={checkForEnter} class="form-control">
                    </div>
                    <div class="col-auto"><span id="pinDialogErrorMessage" class="ms-lg-2 text-danger"></span></div>
                </div>
            </div>

            <div class="modal-footer">
                <span id="pinDialogInfoMessage" class="ms-lg-2 visually-hidden"></span>
                <button type="button" class="btn btn-success btn-sm ms-3" data-bs-dismiss="modal" style="position: absolute; left: 0;" on:click={resetAccount} >Reset Account</button>
                <button id="pinDialogClose" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-primary btn-sm" on:click={setPin}>Login</button>
            </div>
        </div>
    </div>
</div>