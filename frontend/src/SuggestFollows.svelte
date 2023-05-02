<script>
    import {GetContactProfile, FollowContact, RefreshFeedReset} from "../wailsjs/go/main/App.js";

    let follows = [ ]
    const onSuggestFollowsDialog = () => {
        follows = [
            // Some random, popular keys according to primal.net...
            "82341f882b6eabcd2ba7f1ef90aad961cf074af15b9ef44a09f9d2a8fbfbe6a2", // Jack
            "3f770d65d3a764a9c5cb503ae123e62ec7598ad035d836e2a810f3877a745b24", // Derek Ross
            "32e1827635450ebb3c5a7d12c1f8e7b2b514439ac10a67eef3d9fd9c5c68e245", // Will
            "cfe3b4316d905335b6ce056ba0ec230b587a334381e82bf9a02a184f2d068f8d", // Mariezze
            "c48e29f04b482cc01ca1f9ef8c86ef8318c059e0e9353235162f080f26e14c11", // Walker
            "958b754a1d3de5b5eca0fe31d2d555f451325f8498a83da1997b7fcd5c39e88c", // Sersleepy
            "8fb140b4e8ddef97ce4b821d247278a1a4353362623f64021484b372f948000c", // Fishcake
        ];
    }
    window.runtime.EventsOn('evSuggestFollowsDialog', onSuggestFollowsDialog);

    const showError = (msg) => {
        let d = document.getElementById("suggestFollowsErrorMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }
    const showInfo = (msg) => {
        let d = document.getElementById("suggestFollowsInfoMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }

    const removeSuggested = (pk) => {
        let tmp = [];
        for(let a = 0; a < follows.length; a++) {
            let tag = follows[a];
            if(tag !== pk) {
                tmp.push(tag)
            }
            else {
                console.log("Removing " + pk);
            }
        }
        follows = tmp;
    }

    const saveFollows = () => {
        let tmp = [];
        for(let a = 0; a < follows.length; a++) {
            let tag = follows[a];
            if(document.getElementById("sf:" + a).checked) {
                tmp.push(tag)
            }
        }

        FollowContact(tmp).then(()=>{
            RefreshFeedReset("evFollowEventNote");
            document.getElementById("suggestFollowsClose").click();
        }).catch((msg)=>{
            showError(msg);
        });
    }

</script>
<style></style>

<div class="modal" id="suggestFollowsDialog" tabindex="-1" aria-hidden="true" data-bs-backdrop="static" data-bs-keyboard="false">
    <div class="modal-dialog modal-dialog-centered modal-lg" >
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" ><i class="bi-box-arrow-in-right me-3"></i>Follow Suggestions</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <legend>Some suggested follows to get started?</legend>
                Check profiles to follow. You can follow/unfollow these or others later.

                <table class="table mt-3">
                    <tbody>

                    {#if follows.length === 0}
                        <div align="center" class="mx-2"><i>(None)</i></div>
                    {/if}
                    {#each follows as follow, i}
                        {#await GetContactProfile(follow) then p}
                            <tr>
                                <td scope="row"><img src="{p.meta.picture}" alt="" style="width: 36px !important; height: 36px !important; min-width: 36px; min-height: 36px"></td>
                                <td scope="row">{p.meta.display_name || p.meta.name || ""}</td>
                                <td>{p.meta.nip05 || ""}</td>
<!--                                <td><a href="#" on:click={()=>{removeSuggested(p.pk)}}><i class="bi-trash text-muted"/></a></td>-->
                                <td><input class="form-check-input" type="checkbox" value="" id="sf:{i}"></td>
                            </tr>
                        {/await}
                    {/each}

                    </tbody>
                </table>

            </div>
            <div class="modal-footer">
                <label id="getKeysErrorMessage" class="ms-lg-2 text-danger visually-hidden"></label>
                <label id="getKeysInfoMessage" class="ms-lg-2 visually-hidden"></label>
                <button id="suggestFollowsClose" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary btn-sm" on:click={saveFollows} >Save</button>
            </div>
        </div>
    </div>
</div>