import { ExploreWidget } from "components"
import { useState } from "react"
import { AuthRedirectWrapper } from "wrappers"


export const Explore = () => {


  return (
    <AuthRedirectWrapper>
      <ExploreWidget />
    </AuthRedirectWrapper>
  )
}