<script lang="ts">
    import {Notification, notificationStore, NotificationType} from "../../shared/stores/notifications";
    import {fly} from 'svelte/transition';

    let notifications: Notification[] = []

    notificationStore.subscribe(value => {
        if (!value) {
            return
        }
        notifications = [value, ...notifications,]
    })

    function removeNotification(notification: Notification) {
        notifications = notifications.filter(e => e !== notification)
    }

</script>

<style lang="scss">
  @keyframes shrink {
    0% {
      width: 98%;
    }
    100% {
      width: 0;
    }
  }

  .progress {
    position: absolute;
    bottom: 0;
    left: 0;
    background-color: rgb(0, 0, 0, 0.3);
    height: 6px;
    width: 100%;
    animation-name: shrink;
    animation-timing-function: linear;
    animation-fill-mode: forwards;
  }

  .notifications {
    position: absolute;
    top: 0;
    right: 0;
    width: 20vw;
    z-index: 999;
  }
</style>

<ul class="notifications">
    {#each notifications as notification (notification.id)}
        <li>
            <div class="notification is-{notification.type} timeout"
                 in:fly="{{ x: 200, duration: 400 }}"
                 out:fly="{{ x: 200, duration: 500 }}">

                <button class="delete" on:click={() => removeNotification(notification)}></button>
                <strong>{notification.title}</strong>
                <p>{notification.body}</p>
                <div
                        class="progress"
                        style="animation-duration: {notification.timeout}ms;"
                        on:animationend={() => removeNotification(notification) }>
                </div>
            </div>
        </li>
    {/each}
</ul>