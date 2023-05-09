<script>
    /**
     *  Main application component
     *  Most of the backend events are handled here. When the user submits their key, an evPkChange event is
     *  caught, the user icon set, contacts refreshed and the feed loaded.
     */

    import {EventsOn} from "../wailsjs/runtime/runtime.js";
    import Follow from "./Follow.svelte";
    import EventPost from "./EventPost.svelte";
    import {
        Quit,
        GetMyPubkey,
        RefreshContactProfiles,
        GetTextNotesForPubkeys,
        RefreshFeedReset,
        SaveConfigDark,
        GetContactProfile,
        SaveContacts,
        RestoreContacts,
        BeginSubscriptions,
        GetReadableRelays
    } from '../wailsjs/go/main/App.js'
    import { contactStore } from './ContactStore.js'
    import { eventStore, sortedEvents  } from "./EventStore.js";
    import nostrIcon from "./assets/images/nostr.png";
    import loadingGif from "./assets/images/loading.gif";
    import Dialogs from "./Dialogs.svelte";
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";
    import StatusBar from "./StatusBar.svelte";

    let pendingNotes = [];
    let myPk = false;
    let myProfile = false;
    let filtering = false;
    let filterProfile = false;
    let dark = document.documentElement.getAttribute('data-bs-theme') === 'dark';
    let autoRefresh = false;
    let contactPanel = true;

    const onPkChange = (pk) => {
        GetReadableRelays().then((relays)=>{
            if(relays.length === 0) {
                EventsEmit("evMessageDialog", {
                    title: "No Available Relays",
                    message: "No relays able to read from the network. Click OK to configure, or cancel",
                    iconClass: "bi-exclamation-circle",
                    cancelable: true,
                    callback:   ()=>{
                        console.log("Callback...");
                        EventsEmit("evRelayDialog");
                    }
                });
            }
        });

        myPk = pk;
        GetContactProfile(pk).then((p)=>{
            myProfile = p;
            BeginSubscriptions();
        });
    }
    EventsOn('evPkChange', onPkChange);

    const onRefreshNote = (event) => {
        if(autoRefresh) {
            addOrUpdateEvent(event);
        } else {
            pendingNotes.unshift(event);
            pendingNotes = pendingNotes;
        }
    }
    EventsOn('evRefreshNote', onRefreshNote);

    const onFollowEventNote = (event) => {
        addOrUpdateEvent(event);
    }
    EventsOn('evFollowEventNote', onFollowEventNote);

    const addOrUpdateEvent = (event) => {
        let ev = getEventIndex(event);
        if(ev >= 0) {
            eventStore.updateEvent(ev.id, event)
        } else {
            eventStore.addEvent(event);
        }
    }

    const getEventIndex = (event) => {
        for(let a = 0; a < $sortedEvents.length; a++) {
            let c = $sortedEvents[a];
            if(c.id === event.id) {
                return a;
            }
        }
        return -1
    }

    const getDisplayName = (profile) => {
        return profile.meta.display_name || profile.meta.name || profile.meta.nip05 || profile.pk || "";
    }

    const getContactProfileIndex = (profile) => {
        for(let a = 0; a < $contactStore.length; a++) {
            let c = $contactStore[a];
            if(c.pk === profile.pk) {
                return a;
            }
        }
        return -1
    }

    const onMetadata = (profile) => {
        if(!profile.following) {
            return;
        }
        let c = getContactProfileIndex(profile);
        if(c >= 0) {
            $contactStore[c] = profile;
        }
        else {
            $contactStore.push(profile);
        }
    }
    EventsOn('evMetadata', onMetadata);

    const refreshFeed = () => {
        if(pendingNotes.length > 0 && !filtering) {
            for(let a = 0; a < pendingNotes.length; a++) {
                addOrUpdateEvent(pendingNotes[a]);
            }
            pendingNotes = [];
        } else {
            // Force a full refresh
            resetFilterAndRefresh();
        }
    }

    const onRefreshContacts = () => {
        $contactStore = [];
        RefreshContactProfiles();
    }
    EventsOn('evRefreshContacts', onRefreshContacts);

    const launchPostDialog = () => {
        EventsEmit("evPostDialog");
    }
    const launchRelayDialog = () => {
        EventsEmit("evRelayDialog");
    }
    const launchSearchContact = () => {
        EventsEmit("evFindContactDialog");
    }

    const actionQuit = (e) => {
        Quit();
    }

    const launchProfileCard = () => {
        GetMyPubkey().then((pk)=>{
            GetContactProfile(pk).then((p)=>{
                EventsEmit("evProfileCard", p);
            });
        });
    }

    const sortContacts = () => {
        return $contactStore.sort((a, b) => {
            return getDisplayName(b).trim().toLowerCase() < getDisplayName(a).trim().toLowerCase();
        });
    }

    const onFilterByProfile = (profile) => {
        let unixTimeStamp = Math.floor(Date.now() / 1000);
        let week = Math.floor(604800)
        let since = unixTimeStamp - week;
        filtering = true;
        filterProfile = profile;

        eventStore.deleteAll();
        GetTextNotesForPubkeys([profile.pk], "evFollowEventNote", true);
    }
    EventsOn('evFilterByProfile', onFilterByProfile);

    const resetFilterAndRefresh = () => {
        filtering = false;
        pendingNotes = [];
        eventStore.deleteAll();
        RefreshFeedReset();
    }

    const toggleMode = () => {
        dark = !dark;
        if (dark) {
            document.documentElement.setAttribute('data-bs-theme','dark')
        }
        else {
            document.documentElement.setAttribute('data-bs-theme','light')
        }
        SaveConfigDark(dark);
    }
    const toggleRefresh = () => {
        autoRefresh = !autoRefresh;
    }

    const toggleFollows = () => {
        contactPanel = !contactPanel;
        let panelVal = contactPanel ? 'block' : 'none';
        let icnVal = contactPanel ? 'none' : 'inline-block';

        document.getElementById("followsHeader").style.setProperty('display', panelVal);
        document.getElementById("followsPanel").style.setProperty('display', panelVal);
        document.getElementById("contactMenuIcn").style.setProperty('display', icnVal, 'important');
    }

    const saveContacts = () => {
        SaveContacts().then((path)=>{
            EventsEmit("evMessageDialog", {
                title: "Backup Contacts",
                message: "Contacts have been saved to:<br><br><code>" +
                    path + "</code><br><br>If required, use the Restore option to load and re-publish your contacts."
            });
            console.log("Save contacts successful")
        }).catch((e)=>{
            console.error(e);
            EventsEmit("evMessageDialog", {
                title: "Backup Contacts Failed",
                message: "Save contacts failed: " + e,
                iconClass: "bi-exclamation-circle"
            });
        });
    }
    const loadContacts = () => {
        RestoreContacts().then((path)=>{
            EventsEmit("evMessageDialog", {
                title: "Restore Contacts",
                message: "Contacts have been restored and re-published from:<br><br><code>" + path
            });
        }).catch((e)=>{
            console.error(e);
            EventsEmit("evMessageDialog", {
                title: "Backup Contacts Failed",
                message: "Save contacts failed: " + e,
                iconClass: "bi-exclamation-circle"
            });
        });
    }

    $: pendingCount = pendingNotes.length;

</script>

<style></style>

<main>
    <div class="container-fluid d-flex flex-column vh-100 overflow-hidden">
        <nav class="navbar navbar-expand-lg">
            {#if myProfile && myProfile.meta && myProfile.meta.picture}
                <img src="{myProfile.meta.picture}" alt="" height="28" class="me-2">
            {:else}
                <img src="{nostrIcon}" alt="" height="28" class="me-2">
            {/if}
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown">
                        File
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#findEventDialog"><i class="bi bi-search me-3"/>Find Event</a></li>
                        <li><a class="dropdown-item" href="#" on:click={actionQuit}><i class="bi bi-box-arrow-left me-3"/>Quit</a></li>
                    </ul>

                </li>
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown">
                        Post
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#postDialog" on:click={launchPostDialog}><i class="bi bi-file-plus me-3"/>New Post...</a></li>
                        <li><a class="dropdown-item" href="#" on:click={refreshFeed}><i class="bi bi-arrow-clockwise me-3"/>Refresh Feed</a></li>
                    </ul>
                </li>

                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown" >
                        Contact
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#findContactDialog" on:click={launchSearchContact}><i class="bi bi-search me-3"/>Find Contact</a></li>
                        <li><a class="dropdown-item" href="#" on:click={onRefreshContacts}><i class="bi bi-arrow-clockwise me-3"/>Refresh Contact List</a></li>
                        <li>
                            <hr class="dropdown-divider">
                        </li>
                        <li><a class="dropdown-item" href="#"  on:click={saveContacts}><i class="bi bi-file-arrow-down me-3"/>Backup contacts</a></li>
                        <li><a class="dropdown-item" href="#"  on:click={loadContacts}><i class="bi bi-file-arrow-up me-3"/>Restore contacts</a></li>
                    </ul>
                </li>

                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown" >
                        Config
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#relayDialog" on:click={launchRelayDialog}><i class="bi bi-hdd-network me-3"/>Relays</a></li>
                        <li><a class="dropdown-item" href="#" on:click={launchProfileCard}><i class="bi bi-person-badge me-3"/>My Profile</a></li>
                        <li><a class="dropdown-item" href="#loginDialog" data-bs-toggle="modal"><i class="bi-box-arrow-in-right me-3"/>Login</a></li>
                        <li>
                            <hr class="dropdown-divider">
                        </li>
                        <li>
                            <a class="dropdown-item" href="#" >
                            <div class="form-check form-switch" >
                                <input class="form-check-input" on:click={toggleMode} type="checkbox" checked={dark} id="toggleMode">Dark Theme
                            </div>
                            </a>
                        </li>
                        <li>
                            <a class="dropdown-item" href="#" >
                                <div class="form-check form-switch" >
                                    <input class="form-check-input" on:click={toggleRefresh} type="checkbox" checked={autoRefresh} id="toggleRefresh">Auto Refresh
                                </div>
                            </a>
                        </li>
                    </ul>
                </li>
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown" >Help</a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#aboutDialog" data-bs-toggle="modal"  ><i class="bi-info-circle me-3"/>About</a></li>
                    </ul>
                </li>
            </ul>
            <ul class="navbar-nav ms-auto">
                <li class="nav-item "><button class="btn btn-outline-warning me-3" data-bs-toggle="modal" data-bs-target="#postDialog" on:click={launchPostDialog} >Post</button></li>
                <li class="nav-item "><button class="btn btn-outline-success me-3" on:click={refreshFeed}>Refresh
                    {#if pendingCount > 0}
                    <span class="badge bg-success ms-2">{pendingCount}</span>
                    {/if}
                </button></li>
            </ul>
        </nav>
        <div class="row">
            <div class="col-3" id="followsHeader">
                <h6 class="border-bottom pb-2 mb-0"><i class="bi bi-list me-2" style="cursor: pointer;" on:click={toggleFollows} />Follows <span class="ms-2 badge bg-primary-subtle">{$contactStore.length}</span>
                    <span class="float-end">
                            <a href="#" class="text-muted" data-bs-toggle="modal" data-bs-target="#findContactDialog" data-bs-placement="bottom" title="Find contact by key or npub" on:click={launchSearchContact}><i class="bi-search"></i></a>
                            <a href="#" class="text-muted" data-bs-placement="bottom" title="Refresh list" on:click={onRefreshContacts}><i class="bi-arrow-clockwise"></i></a>
                        </span>
                </h6>
            </div>
            <div class="col">
                {#if filtering }
                <h6 id="filter-banner" class="pb-2 mb-0">
                    <a href="#" class="text-muted" title="Back" on:click={resetFilterAndRefresh}>
                        <i class="bi bi-caret-left-fill me-1"></i>
                    </a><span class="text-primary-emphasis"> { getDisplayName(filterProfile) } </span>
                </h6>
                {:else}
                <h6 id="feed-banner" class="pb-2 mb-0"><i id="contactMenuIcn" class="bi bi-list me-2" style="cursor: pointer; display: none;" on:click={toggleFollows} />Subscription Feed</h6>
                {/if}
            </div>
        </div>
        <div class="row flex-grow-1 overflow-hidden">
            <div class="col-3 mh-100 overflow-auto" id="followsPanel">
                <div class="my-3  bg-body">
                    {#each sortContacts($contactStore) as profile}
                        <Follow {profile}/>
                    {/each}
                </div>
            </div>
            <div class='col mh-100 me-2 overflow-auto'>
                <div class='row flex-grow-1 pe-2'>
                    {#if $sortedEvents.length === 0 && myPk}
                        <img src="{loadingGif}" style="width: 50px !important;">
                    {:else}
                        {#each $sortedEvents as event}
                            {#if event.kind === 1 || event.kind === 6 }
                                <EventPost {event} {myPk} />
                            {/if}
                        {/each}
                    {/if}
                </div>
            </div>
        </div>
        <div class="row">
<!--            <div class="col-12" style="background: #0d6efd">-->
            <div class="col-12 border-top mt-2 p-1 px-3">
            <StatusBar />
            </div>
        </div>
    </div>
</main>

<Dialogs />
