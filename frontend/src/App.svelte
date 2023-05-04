<script>
    import Follow from "./Follow.svelte";
    import EventPost from "./EventPost.svelte";
    import Footer from "./Footer.svelte";
    import {
        Quit,
        GetMyPubkey,
        RefreshContactProfiles,
        GetTextNotesByPubkeysOptions,
        RefreshFeedReset,
        SaveConfigDark,
        GetContactProfile
    } from '../wailsjs/go/main/App.js'
    import { eventStore, contactStore } from './store.js'
    import nostrIcon from "./assets/images/nostr.png"
    import defaultProfile from "./Util.svelte"
    import {onMount} from "svelte";
    import Dialogs from "./Dialogs.svelte";

    let pendingNotes = [];
    let myPk;
    let filtering = false;
    let filterProfile = defaultProfile;
    let dark = document.documentElement.getAttribute('data-bs-theme') === 'dark';
    let autoRefresh = false;

    const onRefreshNote = (event) => {
        if(autoRefresh) {
            addOrUpdateEvent(event);
            $eventStore = $eventStore;
        } else {
            pendingNotes.unshift(event);
            pendingNotes = pendingNotes;
        }
    }
    window.runtime.EventsOn('evRefreshNote', onRefreshNote);

    const onFollowEventNote = (event) => {
        addOrUpdateEvent(event);
        $eventStore = $eventStore;
    }
    window.runtime.EventsOn('evFollowEventNote', onFollowEventNote);

    const addOrUpdateEvent = (event) => {
        let c = getEventIndex(event);
        if(c >= 0) {
            $eventStore[c] = event;
        }
        else {
            $eventStore.push(event);
        }
    }

    const getDisplayName = (profile) => {
        return profile.meta.display_name || profile.meta.name || profile.meta.nip05 || profile.pk || "";
    }

    const getEventIndex = (event) => {
        for(let a = 0; a < $eventStore.length; a++) {
            let c = $eventStore[a];
            if(c.id === event.id) {
                return a;
            }
        }
        return -1
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
        $contactStore = $contactStore;
    }
    window.runtime.EventsOn('evMetadata', onMetadata);

    const refreshFeed = () => {
        if(pendingNotes.length > 0) {
            for(let a = 0; a < pendingNotes.length; a++) {
                addOrUpdateEvent(pendingNotes[a]);
            }
            pendingNotes = [];
            $eventStore = $eventStore;
        } else {
            // Force a full refresh
            resetFilterAndRefresh();
        }
    }

    const onRefreshContacts = () => {
        $contactStore = [];
        RefreshContactProfiles();
    }
    window.runtime.EventsOn('evRefreshContacts', onRefreshContacts);

    const newPostDialog = () => {
        window.runtime.EventsEmit("evPostDialog");
    }
    const relayDialog = () => {
        window.runtime.EventsEmit("evRelayDialog");
    }
    const searchContact = () => {
        window.runtime.EventsEmit("evFindContactDialog");
    }

    onMount(async () => {
        const res = await GetMyPubkey();
        myPk = await res;
    });

    const actionQuit = (e) => {
        Quit();
    }

    const myProfile = () => {
        GetMyPubkey().then((pk)=>{
            GetContactProfile(pk).then((p)=>{
                window.runtime.EventsEmit("evProfileCard", p);
            });
        });
    }

    const sortEvents = () => {
        return $eventStore.sort((a, b) => {
            return a.created_at < b.created_at;
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

        $eventStore = [];
        GetTextNotesByPubkeysOptions([profile.pk], [1,6], since, 100, "evFollowEventNote");
    }
    window.runtime.EventsOn('evFilterByProfile', onFilterByProfile);

    const resetFilterAndRefresh = () => {
        filtering = false;
        $eventStore = [];
        RefreshFeedReset("evFollowEventNote");
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


    $: pendingCount = pendingNotes.length;
</script>

<style></style>

<main>
    <div class="container-fluid d-flex flex-column vh-100 overflow-hidden">
        <nav class="navbar navbar-expand-lg">
            <img src={nostrIcon}
                 alt="" height="28" class="me-2">
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown">
                        File
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#findEventDialog"><i class="bi bi-search me-2"/>Find Event</a></li>
                        <li><a class="dropdown-item" href="#" on:click={actionQuit}><i class="bi bi-box-arrow-left me-2"/>Quit</a></li>
                    </ul>
                </li>
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown">
                        Post
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#postDialog" on:click={newPostDialog}><i class="bi bi-file-plus me-2"/>New Post...</a></li>
                        <li><a class="dropdown-item" href="#" on:click={refreshFeed}><i class="bi bi-arrow-clockwise me-2"/>Refresh Feed</a></li>
                    </ul>
                </li>

                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown" >
                        Contact
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#findContactDialog" on:click={searchContact}><i class="bi bi-file-plus me-2"/>Find Contact</a></li>
                        <li><a class="dropdown-item" href="#" on:click={onRefreshContacts}><i class="bi bi-arrow-clockwise me-2"/>Refresh Contact List</a></li>
                        <li>
                            <hr class="dropdown-divider">
                        </li>
                        <li><a class="dropdown-item disabled" href="#"  on:click={()=>{}}><i class="bi bi-file-arrow-down me-2"/>Export to File...</a></li>
                        <li><a class="dropdown-item disabled" href="#"  on:click={()=>{}}><i class="bi bi-file-arrow-up me-2"/>Import from File...</a></li>
                    </ul>
                </li>

                <li class="nav-item dropdown">
                    <a class="nav-link dropdown" href="#" role="button" data-bs-toggle="dropdown" >
                        Config
                    </a>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" data-bs-toggle="modal" data-bs-target="#relayDialog" on:click={relayDialog}><i class="bi bi-hdd-network me-2"/>Relays</a></li>
                        <li><a class="dropdown-item" href="#" on:click={myProfile}><i class="bi bi-person-badge me-2"/>My Profile</a></li>
                        <li><a class="dropdown-item" href="#loginDialog" data-bs-toggle="modal"><i class="bi-box-arrow-in-right me-2"/>Login</a></li>
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
                        <li><a class="dropdown-item" href="#aboutDialog" data-bs-toggle="modal"  ><i class="bi-info-circle me-2"/>About</a></li>
                    </ul>
                </li>
            </ul>
            <ul class="navbar-nav ms-auto">
                <li class="nav-item "><button class="btn btn-outline-warning me-2" data-bs-toggle="modal" data-bs-target="#postDialog" on:click={newPostDialog} >Post</button></li>
                <li class="nav-item "><button class="btn btn-outline-success me-2" on:click={refreshFeed}>Refresh
                    {#if pendingCount > 0}
                    <span class="badge bg-success">{pendingCount}</span>
                    {/if}
                </button></li>
            </ul>
        </nav>
        <div class="row">
            <div class="col-3">
                <h6 class="border-bottom pb-2 mb-0">Follows <span class="ms-2 badge bg-primary-subtle">{$contactStore.length}</span>
                    <span class="float-end">
                            <a href="#" class="text-muted" data-bs-toggle="modal" data-bs-target="#findContactDialog" data-bs-placement="bottom" title="Find contact by key or npub" on:click={searchContact}><i class="bi-search"></i></a>
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
                <h6 id="feed-banner" class="pb-2 mb-0">Subscription Feed</h6>
                {/if}
            </div>
        </div>
        <div class="row flex-grow-1 overflow-hidden">
            <div class="col-3 mh-100 overflow-auto ">
                <div class="my-3  bg-body">
                    {#each sortContacts($contactStore) as profile}
                        <Follow {profile}/>
                    {/each}
                </div>
            </div>
            <div class='col mh-100 me-2 overflow-auto'>
                <div class='row flex-grow-1 pe-2'>
                    {#each sortEvents($eventStore) as event}
                        {#if event.kind === 1 || event.kind === 6 }
                            <EventPost {event} {myPk} />
                        {/if}
                    {/each}
                </div>
            </div>
        </div>
        <div class="row flex-shrink-0 ">
            <Footer/>
        </div>
    </div>
</main>

<Dialogs />
