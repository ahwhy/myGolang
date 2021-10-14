const getters = {
    accessToken: state => state.user.accessToken,
    namespace: state => state.user.namespace,
    account: state => state.user.account,
    username: state => state.user.name,
    userType: state => state.user.type,
    userAvatar: state => state.user.avatar,
  }
  export default getters