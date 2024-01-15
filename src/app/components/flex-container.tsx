import React from "react";

export const FlexContainer = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="w-screen h-screen p-2 overflow-auto md:p-6">
      <div className="w-full h-full max-w-screen-md max-h-screen mx-auto overflow-auto bg-white rounded-xl lg:max-w-800">
        <div className="w-full pb-6">
          <img
            src="/deja-blue-logo.svg"
            alt="Deja Blue logo"
            className="p-6 mx-auto"
          />
          <div className="w-full h-0.5 bg-gray-100  " />
        </div>
        {children}
      </div>
    </div>
  );
};

export default FlexContainer;
